import os
import re

def refactor_usecases():
    directory = r"e:\Bisnis\RekamMedis\backend_go_firebase\internal\usecase"
    for filename in os.listdir(directory):
        if not filename.endswith(".go"): continue
        filepath = os.path.join(directory, filename)
        with open(filepath, 'r') as f:
            content = f.read()

        # Change uint to string in usecases
        # Match signatures like (id uint) -> (id string)
        content = re.sub(r'\(id\s+uint\)', '(id string)', content)
        content = re.sub(r'\(patientID\s+uint\)', '(patientID string)', content)
        content = re.sub(r'\(physioID\s+uint\)', '(physioID string)', content)
        content = re.sub(r'([A-Za-z]+ID)\s+uint\)', r'\1 string)', content)
        content = re.sub(r'([A-Za-z]+ID)\s+uint,', r'\1 string,', content)

        with open(filepath, 'w') as f:
            f.write(content)

def refactor_handlers():
    directory = r"e:\Bisnis\RekamMedis\backend_go_firebase\internal\delivery\http"
    for filename in os.listdir(directory):
        if not filename.endswith(".go"): continue
        filepath = os.path.join(directory, filename)
        with open(filepath, 'r') as f:
            content = f.read()

        # In handlers, parsing logic needs to be removed/changed.
        # Common pattern:
        # idStr := c.Param("id")
        # id, err := strconv.ParseUint(idStr, 10, 32)
        # if err != nil { ... }
        
        # We will use regex to find this block and replace it with:
        # id := c.Param("id")
        
        pattern = re.compile(
            r'idStr\s*:=\s*c\.Param\("([^"]+)"\)\s+'
            r'id(?:64)?,\s*err\s*:=\s*strconv\.ParseUint\(idStr,\s*10,\s*(?:32|64)\)\s+'
            r'if\s+err\s*!=\s*nil\s*\{\s*utils\.ErrorResponse[^\}]+\}\s+'
            r'(?:id\s*:=\s*uint\(id(?:64)?\)\s+)?',
            re.DOTALL
        )
        content = pattern.sub(r'id := c.Param("\1")\n\t', content)

        # sometimes it uses another param name
        pattern2 = re.compile(
            r'([a-zA-Z]+ID)Str\s*:=\s*c\.Param\("([^"]+)"\)\s+'
            r'\1Uint,\s*err\s*:=\s*strconv\.ParseUint\(\1Str,\s*10,\s*32\)\s+'
            r'if\s+err\s*!=\s*nil\s*\{\s*utils\.ErrorResponse[^\}]+\}\s+'
            r'\1\s*:=\s*uint\(\1Uint\)\s+',
            re.DOTALL
        )
        content = pattern2.sub(r'\1 := c.Param("\2")\n\t', content)

        # What if they didn't cast to uint?
        pattern3 = re.compile(
            r'([a-zA-Z]+ID)Str\s*:=\s*c\.Param\("([^"]+)"\)\s+'
            r'([a-zA-Z]+ID),\s*err\s*:=\s*strconv\.ParseUint\([^\)]+\)\s+'
            r'if\s+err\s*!=\s*nil\s*\{\s*utils\.ErrorResponse[^\}]+\}\s+',
            re.DOTALL
        )
        content = pattern3.sub(r'\3 := c.Param("\2")\n\t', content)

        # Also user ID from Context
        # userID, _ := c.Get("user_id")
        # we might need to cast to string if it was uint
        
        # Remove "strconv" import if it's no longer used
        # We'll let `goimports` handle that, or we'll just ignore it (it's fine if unused import exists temporarily, we'll fix it later)

        with open(filepath, 'w') as f:
            f.write(content)

refactor_usecases()
refactor_handlers()
print("Usecases and handlers refactored.")
