import shutil

# NOTE:
# execute this file on backend root directory.

targets = [
    'services/auth/v1/internal/infrastructure/persistence/db/sqlc/',
    'services/media/v1/internal/infrastructure/persistence/db/sqlc/',
]

# delete old directories
for t in targets:
    try:
        shutil.rmtree(t)
    except:
        pass

# copy new directories
for t in targets:
    shutil.copytree('db/sqlc', t)