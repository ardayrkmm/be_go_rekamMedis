import os
import re

directory = r"e:\Bisnis\RekamMedis\backend_go_firebase\internal\repository"

def fix_count():
    for filename in os.listdir(directory):
        if not filename.endswith(".go"): continue
        filepath = os.path.join(directory, filename)
        with open(filepath, 'r') as f:
            content = f.read()
            
        pattern = re.compile(r'q := r\.db\.Collection\([^)]+\)(?:\.Where\([^)]+\))?\s+aggQuery := q\.NewAggregationQuery\(\)\.WithCount\("all"\)\s+res, err := aggQuery\.Get\(ctx\)\s+if err == nil \{\s+total = res\["all"\].*?\s+\}', re.DOTALL)
        content = pattern.sub(r'// total omitted for NoSQL', content)

        with open(filepath, 'w') as f:
            f.write(content)

fix_count()
print("Count fixed.")
