import os
import re

directory = r"e:\Bisnis\RekamMedis\backend_go_firebase\internal\usecase"
for filename in os.listdir(directory):
    if not filename.endswith(".go"): continue
    filepath = os.path.join(directory, filename)
    with open(filepath, 'r') as f:
        content = f.read()

    # Change all `uint` to `string` in usecases (interfaces, method params, variables)
    content = re.sub(r'\buint\b', 'string', content)

    with open(filepath, 'w') as f:
        f.write(content)
        
print("Usecases uint changed to string.")
