import os

dir_path = r'e:\Bisnis\RekamMedis\backend_go_firebase\internal\models'
for f in sorted(os.listdir(dir_path)):
    if not f.endswith('.go'):
        continue
    p = os.path.join(dir_path, f)
    with open(p, 'r', encoding='utf-8') as file:
        c = file.read()
    new_c = c.replace('gorm:"size:36;primaryKey"', 'gorm:"type:varchar(36);primaryKey"')
    new_c = new_c.replace('gorm:"size:36"', 'gorm:"type:varchar(36)"')
    if new_c != c:
        with open(p, 'w', encoding='utf-8') as file:
            file.write(new_c)
        print(f'Fixed: {f}')
print('All done')
