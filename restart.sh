# chmod 744 build
./stop.sh

FILE=build
if test -f "$FILE"; then
    echo "$FILE exists."
    mv "$FILE" main
    chown www-data main
    chgrp www-data main
    chmod 744 main
fi

. ./run.sh

echo 'restart success'
