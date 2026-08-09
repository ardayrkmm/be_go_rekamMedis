import os
import re

dir_path = r'e:\Bisnis\RekamMedis\backend_go_firebase\internal\models'

for filename in os.listdir(dir_path):
    if filename.endswith('.go') and filename != 'setup.go':
        filepath = os.path.join(dir_path, filename)
        with open(filepath, 'r') as f:
            content = f.read()

        # Check if this is a model file
        if 'package models' not in content:
            continue

        # Add import for gorm and uuid if they have a struct with ID string
        if re.search(r'\bID\s+string\b', content):
            if '"gorm.io/gorm"' not in content:
                if 'import (' in content:
                    content = content.replace('import (', 'import (\n\t"github.com/google/uuid"\n\t"gorm.io/gorm"\n', 1)
                else:
                    content = content.replace('package models', 'package models\n\nimport (\n\t"github.com/google/uuid"\n\t"gorm.io/gorm"\n)\n', 1)
                    
            # Add BeforeCreate hook
            structs = re.findall(r'type\s+(\w+)\s+struct\s*\{', content)
            for s in structs:
                # check if struct has ID string
                struct_block_match = re.search(r'type\s+{}\s+struct\s*\{{(.*?)\}}'.format(s), content, re.DOTALL)
                if struct_block_match and re.search(r'\bID\s+string\b', struct_block_match.group(1)):
                    hook = f"""

func (m *{s}) BeforeCreate(tx *gorm.DB) (err error) {{
	if m.ID == "" {{
		m.ID = uuid.New().String()
	}}
	return
}}
"""
                    if hook not in content:
                        content += hook

        # Primary key: ID string `...` -> ID string `... gorm:"type:char(36);primaryKey"`
        def repl(m):
            t = m.group(1)
            if 'gorm:' not in t:
                return f'ID string `{t} gorm:"type:char(36);primaryKey"`'
            return m.group(0)

        content = re.sub(r'ID\s+string\s+`([^`]*)`', repl, content)
        
        with open(filepath, 'w') as f:
            f.write(content)
        print(f"Updated: {filename}")
