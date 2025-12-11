from flask import Flask, request, jsonify
from transformers import AutoTokenizer, AutoModelForSequenceClassification
import logging
import torch # Added import for torch

app = Flask(__name__)
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Load model and tokenizer
MODEL_NAME = "unitary/toxic-bert"
print(f"Loading model: {MODEL_NAME}...")

try:
    tokenizer = AutoTokenizer.from_pretrained(MODEL_NAME)
    model = AutoModelForSequenceClassification.from_pretrained(MODEL_NAME)
    print("Model loaded successfully.")
except Exception as e:
    print(f"Failed to load model: {e}")
    exit(1)

@app.route('/analyze', methods=['POST'])
def analyze():
    data = request.json
    content = data.get('content', '')
    
    if not content:
        return jsonify({'error': 'No content provided'}), 400

    try:
        inputs = tokenizer(content, return_tensors="pt", truncation=True, max_length=512)
        with torch.no_grad():
            outputs = model(**inputs)
        
        # specific to toxic-bert (multilabel classification usually)
        # But unitary/toxic-bert output is typically logits for 6 classes:
        # [toxic, severe_toxic, obscene, threat, insult, identity_hate]
        probs = torch.sigmoid(outputs.logits).squeeze().tolist()
        
        labels = ["toxic", "severe_toxic", "obscene", "threat", "insult", "identity_hate"]
        
        # Build list of flagged categories
        results = []
        is_abusive = False
        max_score = 0.0
        
        for i, score in enumerate(probs):
            if score > 0.5: # Threshold
                results.append({
                    "label": labels[i],
                    "score": score
                })
                is_abusive = True
                if score > max_score:
                    max_score = score
        
        # Log safely
        content_preview = content[:20] if content else ""
        print(f"Analyzed: '{content_preview}...' -> Flags: {results}")

        return jsonify({
            "is_abusive": is_abusive,
            "confidence_score": max_score,
            "flags": results
        })

    except Exception as e:
        logger.error(f"Inference error: {e}")
        return jsonify({"error": "Analysis failed"}), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5001)
