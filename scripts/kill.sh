LSOF=$(lsof -ti :$1)
if [[ $LSOF ]]; then
  kill -9 $LSOF
fi
