runit() {
  echo $FCL
  OUTDIR=${NAME}_${PROC}
  if [ -z "$OUTDIR" ]; then
    echo Output directory not specified.
    return 1
  fi
  if [ -r $OUTDIR ]; then
    echo Output directory already exists.
    echo Please remove $OUTDIR
    return 2
  fi
  ARGS="-c $FCL"
  if [ -n "$INFILE" ]; then
    ARGS="$ARGS -s ../$INFILE"
  fi
  if [ -n "$OUTFILE" ]; then
    ARGS="$ARGS -o $OUTFILE"
  fi
  EVENT_COUNT=$NEVT
  if [ -z "$INFILE" -a -z "$NEVT" ]; then
    EVENT_COUNT=1000
  fi
  if [ -n "$EVENT_COUNT" ]; then
    ARGS="$ARGS -n $EVENT_COUNT"
  fi
  if [ -n "$NSKP" ]; then
    ARGS="$ARGS -n $NSKP"
  fi
  echo "Args: $ARGS"
  mkdir $OUTDIR
  cd $OUTDIR
  FHICL_FILE_PATH="../fcl:$FHICL_FILE_PATH"
  fcldump $FCL 9 >${FCL}dump
  lar $ARGS 2>&1 | tee run.log
  STAT=$?
  echo Output directory: $OUTDIR
  echo Return Status: $STAT
  return $STAT
}

NAME=$1
PROC=$2
NEVT=$3
NSKP=$4
if [ -z "$PROC" -o -z "$NAME" -o "$PROC" = "-h" ]; then
  echo Usage $0: NAME PROC
  echo "  NAME = mu1000, ..."
  echo "  PROC = gen, g4, detsim, dataprep, ..."
  exit 0
fi

FCL=$PROC.fcl
INPROC=
OUTFILE=${NAME}_${PROC}_event.root
if [ $PROC = "gen" ]; then
  FCL=${PROC}_${NAME}.fcl
elif [ $PROC = g4 ]; then
  INPROC=gen
elif [ $PROC = detsim ]; then
  INPROC=g4
elif [ $PROC = dataprep ]; then
  INPROC=detsim
elif [ $PROC = dataprep-old ]; then
  INPROC=detsim
else
  echo "Unknown process: $PROC"
  exit 3
fi
FCLFILE=fcl/$FCL
if [ ! -r $FCLFILE ]; then
  echo FCL not found: $FCLFILE
  exit 1
fi
INFILE=
if [ -n "$INPROC" ]; then
  INFILE=${NAME}_${INPROC}/${NAME}_${INPROC}_event.root
  if [ ! -r $INFILE ]; then
    echo Input event file not found: $INFILE
    exit 4
  fi
fi 

runit
