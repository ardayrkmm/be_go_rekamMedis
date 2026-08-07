import os
import re

directory = r"e:\Bisnis\RekamMedis\backend_go\internal\delivery\http"

for filename in os.listdir(directory):
    if not filename.endswith(".go"):
        continue
    filepath = os.path.join(directory, filename)
    with open(filepath, 'r') as f:
        content = f.read()

    # Replace FormatValidationError
    content = re.sub(
        r'c\.JSON\(http\.StatusUnprocessableEntity,\s*utils\.FormatValidationError\(err\)\)',
        r'utils.HandleValidationError(c, err)',
        content
    )

    # Replace gin.H{"error": ...} -> utils.ErrorResponse
    # Match: c.JSON(http.Status..., gin.H{"error": "..."})
    content = re.sub(
        r'c\.JSON\((http\.[A-Za-z0-9_]+),\s*gin\.H\{"error":\s*([^}]+)\}\)',
        r'utils.ErrorResponse(c, \1, \2, nil)',
        content
    )

    # Replace gin.H{"message": ...} -> utils.SuccessResponse(..., nil)
    # Match: c.JSON(http.Status..., gin.H{"message": "..."})
    content = re.sub(
        r'c\.JSON\((http\.[A-Za-z0-9_]+),\s*gin\.H\{"message":\s*([^}]+)\}\)',
        r'utils.SuccessResponse(c, \1, \2, nil)',
        content
    )

    # Replace gin.H{"data": ...} -> utils.SuccessResponse(..., "Success", ...)
    content = re.sub(
        r'c\.JSON\((http\.[A-Za-z0-9_]+),\s*gin\.H\{"data":\s*([^}]+)\}\)',
        r'utils.SuccessResponse(c, \1, "Success", \2)',
        content
    )

    # Replace gin.H{"data": ..., "token": ...} -> utils.SuccessResponse(..., "Success", gin.H{...})
    # This is for AuthHandler (we'll just manually fix AuthHandler if needed, let's see if it gets mangled)

    # Replace pagination
    # Match: c.JSON(http.StatusOK, gin.H{ "data": ..., "total": ..., "page": ... })
    content = re.sub(
        r'c\.JSON\((http\.[A-Za-z0-9_]+),\s*gin\.H\{\s*"data":\s*([^,]+),\s*"total":\s*([^,]+),\s*"page":\s*([^,]+),\s*\}\)',
        r'utils.SuccessResponsePaginated(c, \1, "Success", \2, \4, perPage, \3)',
        content
    )

    with open(filepath, 'w') as f:
        f.write(content)
