# python hf Tiny Model Service
#
# SmolLM2-135M-Instruct from hugging face
#
import os
import logging
from contextlib import asynccontextmanager
from typing import List, Optional

import torch
from fastapi import FastAPI, HTTPException, status
from pydantic import BaseModel, Field
from transformers import AutoModelForCausalLM, AutoTokenizer, TextIteratorStreamer
from threading import Thread

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("inference-py")

MODEL_ID = os.getenv("MODEL_ID", "HuggingFaceTB/SmolLM2-135M-Instruct")
DEVICE = "cuda" if torch.cuda.is_available() else "cpu"

# Global references for model and tokenizer
model = None
tokenizer = None


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Load model and tokenizer during startup."""
    global model, tokenizer
    logger.info(f"Loading model '{MODEL_ID}' on device: {DEVICE}...")

    try:
        tokenizer = AutoTokenizer.from_pretrained(MODEL_ID)
        model = AutoModelForCausalLM.from_pretrained(
            MODEL_ID,
            torch_dtype=torch.float16 if DEVICE == "cuda" else torch.float32,
            device_map="auto" if DEVICE == "cuda" else None,
        )
        if DEVICE == "cpu":
            model.to("cpu")
        logger.info("Model loaded successfully.")
    except Exception as e:
        logger.error(f"Failed to load model: {e}")
        raise e

    yield

    # Cleanup logic
    logger.info("Shutting down inference service...")


app = FastAPI(title="Inference Microservice", version="1.0.0", lifespan=lifespan)


# Request/Response Schemas
class ChatMessage(BaseModel):
    role: str = Field(..., description="'user', 'assistant', or 'system'")
    content: str


class ChatCompletionRequest(BaseModel):
    messages: List[ChatMessage]
    max_new_tokens: Optional[int] = Field(default=256, ge=1, le=1024)
    temperature: Optional[float] = Field(default=0.7, ge=0.0, le=2.0)
    top_p: Optional[float] = Field(default=0.9, ge=0.0, le=1.0)


class ChatCompletionResponse(BaseModel):
    response: str
    model: str


# Endpoints
@app.get("/healthz", status_code=status.HTTP_200_OK)
async def health_check():
    """Liveness probe endpoint for Kubernetes."""
    if model is None or tokenizer is None:
        raise HTTPException(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE, detail="Model not ready"
        )
    return {"status": "ok", "model": MODEL_ID, "device": DEVICE}


@app.post("/v1/chat/completions", response_model=ChatCompletionResponse)
async def generate_chat_completion(request: ChatCompletionRequest):
    """Standard non-streaming chat generation endpoint."""
    if model is None or tokenizer is None:
        raise HTTPException(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
            detail="Model is not initialized.",
        )

    try:
        # Convert Pydantic objects to dict for tokenizer template
        formatted_messages = [msg.model_dump() for msg in request.messages]

        # Apply standard chat template
        input_text = tokenizer.apply_chat_template(
            formatted_messages, tokenize=False, add_generation_prompt=True
        )

        inputs = tokenizer(input_text, return_tensors="pt").to(DEVICE)

        with torch.no_grad():
            outputs = model.generate(
                **inputs,
                max_new_tokens=request.max_new_tokens,
                temperature=request.temperature if request.temperature > 0 else 1.0,
                top_p=request.top_p,
                do_sample=request.temperature > 0,
                pad_token_id=tokenizer.eos_token_id,
            )

        # Decode response excluding the prompt inputs
        generated_tokens = outputs[0][inputs["input_ids"].shape[1] :]
        response_text = tokenizer.decode(generated_tokens, skip_special_tokens=True)

        return ChatCompletionResponse(response=response_text, model=MODEL_ID)

    except Exception as e:
        logger.error(f"Inference error: {e}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Inference generation failed: {str(e)}",
        )


if __name__ == "__main__":
    import uvicorn

    uvicorn.run("main:app", host="0.0.0.0", port=8000, reload=False)
