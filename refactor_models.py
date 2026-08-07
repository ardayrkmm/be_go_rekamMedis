import os
import re

directory = r"e:\Bisnis\RekamMedis\backend_go_firebase\internal\models"

for filename in os.listdir(directory):
    if not filename.endswith(".go"):
        continue
    filepath = os.path.join(directory, filename)
    with open(filepath, 'r') as f:
        content = f.read()

    # 1. Remove "gorm.io/gorm" import
    content = re.sub(r'\s*"gorm\.io/gorm"', '', content)

    # 2. Change `ID uint` to `ID string`
    content = re.sub(r'ID\s+uint\s+`gorm:"primaryKey"[^`]*`', r'ID string `firestore:"id,omitempty" json:"id"`', content)
    
    # 3. Change other `uint` fields to `string` (usually foreign keys)
    content = re.sub(r'([A-Za-z]+ID)\s+uint\s+`', r'\1 string `', content)

    # 4. Remove `gorm:"..."` tags but keep `json:"..."`
    # Match: `gorm:"..." json:"..."` -> `json:"..."`
    content = re.sub(r'`gorm:"[^"]*"\s+(json:"[^"]*")`', r'`\1`', content)
    
    # 5. Sometimes json is first `json:"..." gorm:"..."`
    content = re.sub(r'`(json:"[^"]*")\s+gorm:"[^"]*"`', r'`\1`', content)
    
    # 6. If only gorm is present: `gorm:"..."` -> remove it completely or just backticks
    content = re.sub(r'`gorm:"[^"]*"`', '', content)
    
    # 7. Remove empty backticks
    content = re.sub(r'\s*``', '', content)
    
    # 8. Change DeletedAt gorm.DeletedAt to DeletedAt *time.Time
    if 'DeletedAt' in content and 'gorm.DeletedAt' in content:
        content = re.sub(r'DeletedAt\s+gorm\.DeletedAt[^`]*`[^`]*`', r'DeletedAt *time.Time `json:"-" firestore:"deleted_at,omitempty"`', content)
        content = re.sub(r'DeletedAt\s+gorm\.DeletedAt.*', r'DeletedAt *time.Time `json:"-" firestore:"deleted_at,omitempty"`', content)
        
        # Add import "time" if not exists
        if '"time"' not in content:
            if 'import (' in content:
                content = content.replace('import (', 'import (\n\t"time"')
            else:
                content = content.replace('package models\n', 'package models\n\nimport "time"\n')

    # Remove extra blank lines caused by import removal
    content = re.sub(r'\n{3,}', '\n\n', content)

    with open(filepath, 'w') as f:
        f.write(content)

print("Refactored models.")
