import os
import re

directory = r"e:\Bisnis\RekamMedis\backend_go_firebase\internal\delivery\http"

for filename in os.listdir(directory):
    if not filename.endswith(".go"): continue
    filepath = os.path.join(directory, filename)
    with open(filepath, 'r') as f:
        content = f.read()

    # Replace id parsing blocks
    # Match:
    # idStr := c.Param("id")
    # id, err := strconv.ParseUint(idStr, 10, 32)
    # if err != nil { ... }
    # 
    # Or:
    # id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    
    # Generic replace:
    content = re.sub(
        r'idStr\s*:=\s*c\.Param\("([^"]+)"\)\s+id(?:64)?,\s*err\s*:=\s*strconv\.ParseUint\([^)]+\)\s+if\s+err\s*!=\s*nil\s*\{[^}]+\}\s+',
        r'id := c.Param("\1")\n\t',
        content, flags=re.DOTALL
    )
    
    # Same for patientIDStr, physioIDStr, etc.
    content = re.sub(
        r'([a-zA-Z]+ID)Str\s*:=\s*c\.Param\("([^"]+)"\)\s+\1(?:64|Uint)?,\s*err\s*:=\s*strconv\.ParseUint\([^)]+\)\s+if\s+err\s*!=\s*nil\s*\{[^}]+\}\s+(?:\1\s*:=\s*uint\(\1(?:64|Uint)\)\s+)?',
        r'\1 := c.Param("\2")\n\t',
        content, flags=re.DOTALL
    )

    # Some might do: id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    content = re.sub(
        r'id(?:64)?,\s*err\s*:=\s*strconv\.ParseUint\(c\.Param\("([^"]+)"\),\s*10,\s*32\)\s+if\s+err\s*!=\s*nil\s*\{[^}]+\}\s+(?:id\s*:=\s*uint\(id(?:64)?\)\s+)?',
        r'id := c.Param("\1")\n\t',
        content, flags=re.DOTALL
    )
    
    # Any residual uint() casts
    content = re.sub(r'uint\((id)\)', r'\1', content)

    with open(filepath, 'w') as f:
        f.write(content)
