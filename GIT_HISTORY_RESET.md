# Git History Reset Instructions

To delete all commit history and start fresh with v1.0.0:

1. Create a new orphan branch (no history):
   git checkout --orphan new-main

2. Add all files:
   git add -A

3. Commit as v1.0.0:
   git commit -m "Initial commit - v1.0.0"

4. Delete the old main branch:
   git branch -D main

5. Rename the new branch to main:
   git branch -m main

6. Force push to remote (WARNING: This deletes all history):
   git push -f origin main

7. Update remote repository name on GitHub:
   - Go to repository settings
   - Rename from "forge" to "cpx"

8. Update remote URL:
   git remote set-url origin https://github.com/ozacod/cpx.git

9. Tag the initial release:
   git tag v1.0.0
   git push origin v1.0.0
