import os
import re

directory = r"e:\Bisnis\RekamMedis\backend_go_firebase\internal\delivery\http"

for filename in os.listdir(directory):
    if not filename.endswith(".go"): continue
    filepath = os.path.join(directory, filename)
    with open(filepath, 'r') as f:
        content = f.read()

    # Replace uint(id) with id
    content = re.sub(r'uint\((id(?:[a-zA-Z0-9_]+)?)\)', r'\1', content)
    
    # Replace uint(userID.(float64)) with userID.(string)
    content = re.sub(r'uint\(([a-zA-Z]+)\.\(float64\)\)', r'\1.(string)', content)

    # Some variables like `userID.(float64)` directly might be used, replace `.(float64)` if it's user_id
    # But only in the context of IDs
    content = re.sub(r'([a-zA-Z]*ID)\.\(float64\)', r'\1.(string)', content)

    with open(filepath, 'w') as f:
        f.write(content)
