import os
import re

directory = r"e:\Bisnis\RekamMedis\backend_go_firebase\internal\delivery\http"

for filename in os.listdir(directory):
    if not filename.endswith(".go"): continue
    filepath = os.path.join(directory, filename)
    with open(filepath, 'r') as f:
        content = f.read()

    # Generic fix for any ParseUint(c.Param("id"), ...) left
    # Match:
    # patientID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    # if err != nil {
    # 	utils.ErrorResponse(...)
    # 	return
    # }
    
    pattern = re.compile(
        r'([a-zA-Z]+ID(?:[a-zA-Z0-9_]*)?),\s*err\s*:=\s*strconv\.ParseUint\(c\.Param\("([^"]+)"\),\s*10,\s*32\)\s+'
        r'if\s+err\s*!=\s*nil\s*\{[^}]+\}\s+',
        re.DOTALL
    )
    content = pattern.sub(r'\1 := c.Param("\2")\n\t', content)

    # Any remaining ParseUint
    pattern2 = re.compile(
        r'([a-zA-Z]+ID(?:[a-zA-Z0-9_]*)?),\s*err\s*:=\s*strconv\.ParseUint\(([^)]+)\,\s*10,\s*32\)\s+'
        r'if\s+err\s*!=\s*nil\s*\{[^}]+\}\s+',
        re.DOTALL
    )
    content = pattern2.sub(r'\1 := \2\n\t', content)
    
    with open(filepath, 'w') as f:
        f.write(content)
