import os
import re

src_dir = r'e:\Bisnis\RekamMedis\backend_go\internal\models'
dst_dir = r'e:\Bisnis\RekamMedis\backend_go_firebase\internal\models'

if not os.path.exists(dst_dir):
    os.makedirs(dst_dir)

for filename in os.listdir(src_dir):
    if filename.endswith('.go'):
        src_file = os.path.join(src_dir, filename)
        dst_file = os.path.join(dst_dir, filename)
        
        with open(src_file, 'r') as f:
            content = f.read()
            
        # Replace uint fields with string
        content = re.sub(r'\bID\s+uint\b', 'ID string', content)
        content = re.sub(r'(\w+ID)\s+\*?uint\b', r'\1 string', content)
        
        # Replace gorm:"primaryKey" with gorm:"type:char(36);primaryKey"
        content = re.sub(r'gorm:"primaryKey"', r'gorm:"type:char(36);primaryKey"', content)
        
        # We also need to remove the ',string' from json tags because they are already string!
        content = re.sub(r'json:"([^"]+),string"', r'json:"\1"', content)

        # Add uuid import if we add BeforeCreate
        struct_match = re.search(r'type (\w+) struct \{', content)
        if struct_match and filename not in ['setup.go']:
            struct_name = struct_match.group(1)
            
            # Check if it has an ID field
            if 'ID string' in content:
                # Add import "github.com/google/uuid" if not present
                if '"github.com/google/uuid"' not in content:
                    content = content.replace('import (', 'import (\n\t"github.com/google/uuid"\n', 1)
                    if 'import (' not in content: # single import fallback
                        content = re.sub(r'import\s+"([^"]+)"', r'import (\n\t"\1"\n\t"github.com/google/uuid"\n)', content, 1)

                hook = f"""

func (m *{struct_name}) BeforeCreate(tx *gorm.DB) (err error) {{
	if m.ID == "" {{
		m.ID = uuid.New().String()
	}}
	return
}}
"""
                content += hook

        with open(dst_file, 'w') as f:
            f.write(content)
        print(f"Migrated model: {filename}")
