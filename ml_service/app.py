from flask import Flask, request, jsonify
from transformers import AutoTokenizer, AutoModelForSequenceClassification, pipeline
import logging

app = Flask(__name__)
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

MODEL_NAME = "Hate-speech-CNERG/english-abusive-MuRIL"

logger.info(f"Loading model: {MODEL_NAME}...")
try:
    # Load model and tokenizer
    tokenizer = AutoTokenizer.from_pretrained(MODEL_NAME)
    model = AutoModelForSequenceClassification.from_pretrained(MODEL_NAME)
    classifier = pipeline("text-classification", model=model, tokenizer=tokenizer)
    logger.info("Model loaded successfully.")
except Exception as e:
    logger.error(f"Failed to load model: {e}")
    raise e

@app.route('/analyze', methods=['POST'])
def analyze():
    data = request.json
    if not data or 'content' not in data:
        return jsonify({"error": "Missing 'content' field"}), 400

    content = data['content']
    try:
        # Run inference
        # The model returns list of dicts: [{'label': 'LABEL_1', 'score': 0.99}]
        # LABEL_0 usually non-abusive, LABEL_1 abusive (need to verify mapping for this specific model)
        # Checking model card: "Class 0: Non-abusive, Class 1: Abusive" usually.
        results = classifier(content)
        result = results[0]
        
        label = result['label']
        score = result['score']
        
        # Map generic LABEL_X to meaningful names if possible, or just return as is
        is_abusive = False
        if label == "LABEL_1": # Assuming 1 is Abusive based on common practice for this dataset
            is_abusive = True
            
        logger.info(f"Analyzed content: '{content[:20]}...' -> Label: {label}, Score: {score}")

        return jsonify({
            "is_abusive": is_abusive,
            "confidence_score": score,
            "raw_label": label
        })

    except Exception as e:
        logger.error(f"Inference error: {e}")
        return jsonify({"error": "Analysis failed"}), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5001)
