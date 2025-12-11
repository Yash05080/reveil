# Révéil API 🌙

> **A Secure, Private, and Safe Social API Platform**
> Built with Go, Python, and Advanced AI Moderation.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![Python](https://img.shields.io/badge/Python-3.9+-3776AB?style=for-the-badge&logo=python)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-4169E1?style=for-the-badge&logo=postgresql)
![Status](https://img.shields.io/badge/Status-Production%20Ready-success?style=for-the-badge)

## 📖 Overview

**Révéil API** is a next-generation backend for private community building. It prioritizes user safety and privacy through rigorous end-to-end encryption and a multi-stage AI moderation pipeline.

Designed for high performance and scalability, it features real-time social interactions, nested commenting, and a comprehensive safety suite including user reporting and automated toxicity analysis.

## ✨ Key Features

### 🔒 Privacy First
- **AES-256-GCM Encryption**: All user content (posts & titles) is encrypted at rest using community-scoped keys.
- **Client SDKs**: Ready-to-use encryption libraries for **Web, iOS/Swift, Android/Kotlin, Flutter**, and **Python**.

### 🛡️ Advanced Safety & Moderation
- **Dual-Layer Pipeline**:
    1.  **Light Layer (Go)**: Instant blocking of known hate speech and restricted keywords (Suicide/Self-harm detection).
    2.  **Heavy Layer (Python/ML)**: Asynchronous deep analysis using **BERT** models to detect subtle toxicity, threats, and insults.
- **Re-Moderation**: Edited content is automatically re-scanned to prevent evasion.
- **User Reporting**: Integrated flagging system for community-led safety.

### ⚡ Performance & Real-time
- **Rate Limiting**: Built-in Token Bucket rate limiter (5000 reqs/sec burst capacity) to prevent abuse.
- **SSE (Server-Sent Events)**: Real-time updates for new posts and moderation actions.
- **Optimized SQL**: Efficient pagination and indexing on Supabase (PostgreSQL).

## 🤖 AI Models

We utilize state-of-the-art Natural Language Processing (NLP) for safety:

| Model | Purpose | Performance |
|-------|---------|-------------|
| **`unitary/toxic-bert`** | Primary toxicity classification (6 labels) | ~92% Accuracy |
| **Keyword/Regex** | Instant blocking of severe threats | <1ms Latency |

*Labels detected: `toxic`, `severe_toxic`, `obscene`, `threat`, `insult`, `identity_hate`.*

## 🛠️ Tech Stack

- **Backend**: Golang (Standard Library + Gorilla Mux)
- **ML Service**: Python (Flask + Transformers + PyTorch)
- **Database**: PostgreSQL (via Supabase)
- **Authentication**: JWT (Supabase Auth)
- **Infrastructure**: Docker Ready

## 🚀 Getting Started

### Prerequisites
- Go 1.21+
- Python 3.9+
- Supabase Project

### Installation

1.  **Clone the repository**
    ```bash
    git clone https://github.com/yourusername/reveil-api.git
    cd reveil-api
    ```

2.  **Environment Setup**
    Create a `.env` file (NOT committed to git):
    ```env
    # Server
    PORT=8080
    LOG_LEVEL=info
    MASTER_ENCRYPTION_KEY=your_32_byte_master_key_base64

    # Supabase
    SUPABASE_URL=https://your-project.supabase.co
    SUPABASE_SERVICE_KEY=your_service_role_key
    JWT_SECRET=your_jwt_secret
    ```

3.  **Run the ML Service**
    ```bash
    cd ml_service
    pip install -r requirements.txt
    python app.py
    ```

4.  **Run the API Server**
    ```bash
    # In root directory
    go mod tidy
    go run main.go
    ```

## 📦 Client Integration (SDKs)

We provide official encryption utilities for seamless frontend integration. Find them in the `client_libs/` directory:

- **Web (JS)**: `client_libs/js/encryption.js`
- **TypeScript**: `client_libs/ts/encryption.ts`
- **Android (Kotlin)**: `client_libs/encryption_util.kt`
- **Flutter (Dart)**: `client_libs/dart/encryption.dart`
- **Python**: `client_libs/python/encryption.py`

*See [Frontend Integration Guide](docs/frontend_integration_guide.md) for full API documentation.*

## 🔒 Security

- **Secrets Management**: `.gitignore` is pre-configured to exclude all sensitive keys.
- **Vulnerability Scanning**: Dependencies are kept minimal to reduce attack surface.
- **Audit Logs**: Moderation actions are permanently logged in `moderation_flags`.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<p align="center">
  Built with ❤️ for a safer internet.
</p>
