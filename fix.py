import os

files = ['activity_log.go', 'gender.go', 'notification.go', 'patient_category.go']
dir_path = r'e:\Bisnis\RekamMedis\backend_go_firebase\internal\models'

for f in files:
    filepath = os.path.join(dir_path, f)
    with open(filepath, 'r') as file:
        content = file.read()
    if '"gorm.io/gorm"' not in content:
        if 'import (' in content:
            content = content.replace('import (', 'import (\n\t"gorm.io/gorm"\n', 1)
        else:
            content = content.replace('package models', 'package models\n\nimport "gorm.io/gorm"\n', 1)
    with open(filepath, 'w') as file:
        file.write(content)
