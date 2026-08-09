import os, glob, re
for path in glob.glob('internal/repository/*.go'):
    with open(path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    new_content = re.sub(r'First\(&([a-zA-Z]+), id\)', r'First(&\1, "id = ?", id)', content)
    
    if new_content != content:
        with open(path, 'w', encoding='utf-8') as f:
            f.write(new_content)
        print('Updated', path)
