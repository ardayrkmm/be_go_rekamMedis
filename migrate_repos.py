import os
import re

src_dir = r'e:\Bisnis\RekamMedis\backend_go\internal\repository'
dst_dir = r'e:\Bisnis\RekamMedis\backend_go_firebase\internal\repository'

if not os.path.exists(dst_dir):
    os.makedirs(dst_dir)

for filename in os.listdir(src_dir):
    if filename.endswith('.go'):
        src_file = os.path.join(src_dir, filename)
        dst_file = os.path.join(dst_dir, filename)
        
        with open(src_file, 'r') as f:
            content = f.read()
            
        # Replace uint with string for IDs
        content = re.sub(r'\b(id\s+)uint\b', r'\1string', content)
        content = re.sub(r'\b(patientID\s+)uint\b', r'\1string', content)
        content = re.sub(r'\b(physioID\s+)uint\b', r'\1string', content)
        content = re.sub(r'\b(categoryID\s+)uint\b', r'\1string', content)
        content = re.sub(r'\b(userID\s+)uint\b', r'\1string', content)
        
        # In Go, FindByID(id string) is fine with GORM.
        
        with open(dst_file, 'w') as f:
            f.write(content)
        print(f"Migrated repository: {filename}")
