import os
import re

model_dir = r'internal\models'
for filename in os.listdir(model_dir):
    if not filename.endswith('.go'):
        continue
    path = os.path.join(model_dir, filename)
    with open(path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Fix ServiceID -> ServiceMasterID
    content = content.replace('foreignKey:ServiceID', 'foreignKey:ServiceMasterID')
    
    # Add constraint:- if missing
    def repl(m):
        inner = m.group(1)
        if 'constraint:-' not in inner:
            if 'foreignKey:' in inner:
                inner = inner + ';constraint:-'
        return f'gorm:"{inner}"'
        
    content = re.sub(r'gorm:\"(foreignKey:[^\"]+)\"', repl, content)
    
    with open(path, 'w', encoding='utf-8') as f:
        f.write(content)
    print(f'Fixed {filename}')
