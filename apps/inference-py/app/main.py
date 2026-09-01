import os
import logging
from contextlib import asynccontextmanager
from typing import List, Optional

from fastapi import FastAPI, HTTPException, status
from llama_cpp import Llama
from pydantic import BaseModel, Field

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("inference-py")

# Download/store a GGUF format model locally or via GGUF repo
# e.g., HuggingFaceTB/SmolLM2-135M-Instruct-GGUF
MODEL_PATH = os.getenv("MODEL_PATH", "./models/smollm2-135m-instruct-q4_k_m.gguf")

llm: Optional[Llama] = None


@asynccontextmanager
async def lifespan(app: FastAPI):
    global llm
    logger.info("Loading quantized GGUF model...")
    try:
        llm = Llama(
            model_path=MODEL_PATH,
            n_ctx=2048,
            n_threads=1,  # Match single vCPU on t2.micro
        )
        logger.info("GGUF Model loaded successfully.")
    except Exception as e:
        logger.error(f"Failed to load model: {e}")
        raise e
    yield


app = FastAPI(title="Inference Service", lifespan=lifespan)


class ChatMessage(BaseModel):
    role: str
    content: str


class ChatCompletionRequest(BaseModel):
    messages: List[ChatMessage]
    max_new_tokens: Optional[int] = Field(default=256, ge=1, le=1024)
    temperature: Optional[float] = Field(default=0.7, ge=0.0, le=2.0)


class ChatCompletionResponse(BaseModel):
    response: str


@app.get("/healthz")
async def health_check():
    if llm is None:
        raise HTTPException(status_code=503, detail="Model not ready")
    return {"status": "ok"}


@app.post("/v1/chat/completions", response_model=ChatCompletionResponse)
async def generate_chat_completion(request: ChatCompletionRequest):
    if llm is None:
        raise HTTPException(status_code=503, detail="Model uninitialized")

    try:
        # Format input for OpenAI-compatible chat format supported by llama-cpp
        formatted_messages = [msg.model_dump() for msg in request.messages]

        response = llm.create_chat_completion(
            messages=formatted_messages,
            max_tokens=request.max_new_tokens,
            temperature=request.temperature,
        )

        output_text = response["choices"][0]["message"]["content"]
        return ChatCompletionResponse(response=output_text)

    except Exception as e:
        logger.error(f"Inference error: {e}")
        raise HTTPException(status_code=500, detail=str(e))
