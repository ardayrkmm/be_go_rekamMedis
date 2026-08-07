import glob
import os

for f in glob.glob('internal/models/*.go'):
    with open(f, 'r') as file:
        content = file.read()
    
    content = content.replace('"deleted_at,omitempty"', '"DeletedAt"')
    
    with open(f, 'w') as file:
        file.write(content)
print("Done")
