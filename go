#!/bin/sh

runit() {
  echo $FCL
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
    ARGS="$ARGS -s $INFILE"
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
  FHICL_FILE_PATH="`pwd`/fcl:$FHICL_FILE_PATH"
  mkdir $OUTDIR
  cd $OUTDIR
  fcldump $FCL 9 >${FCL}dump
  lar $ARGS 2>&1 | tee run.log
  STAT=$?
  echo Output directory: $OUTDIR
  echo Return Status: $STAT
  return $STAT
}

NAME=$1
CPROC=$2
NEVT=$3
NSKP=$4
if [ -z "$CPROC" -o -z "$NAME" -o "$CPROC" = "-h" ]; then
  echo Usage $0: NAME CPROC NEVT NSKP
  echo "  NAME = mu1000, ..."
  echo "  CPROC = PROC or INPROC/PROC"
  echo "  PROC = gen, g4, detsim, dataprep, ..."
  exit 0
fi

PROC=`basename $CPROC`
INPROC=`dirname $CPROC`
if [ $INPROC = "." ]; then 
  if [ $PROC = "gen" ]; then
    FCL=${PROC}_${NAME}.fcl
  elif [ $PROC = g4 ]; then
    INPROC=gen
  elif [ ${PROC:0:6} = detsim ]; then
    INPROC=g4
  elif [ ${PROC:0:8} = dataprep ]; then
    INPROC=detsim
  elif [ $PROC = dataprep-old ]; then
    INPROC=detsim
  elif [ $PROC = dataprep-skip ]; then
    INPROC=detsim
  elif [ $PROC = dataprep-keeprois ]; then
    INPROC=detsim
  elif [ $PROC = chanmap ]; then
    INPROC=detsim
  else
    echo "Unknown process: $PROC"
    exit 3
  fi
  OUTDIR=${NAME}_${PROC}
else
  INDIR=..
  OUTDIR=${NAME}_$INPROC/${PROC}
fi
FCL=$PROC.fcl
OUTFILE=${NAME}_${PROC}_event.root
FCLFILE=fcl/$FCL
if [ ! -r $FCLFILE ]; then
  echo FCL not found: $FCLFILE
  exit 1
fi
INFILE=
if [ -n "$INPROC" ]; then
  INFILE=`pwd`/${NAME}_${INPROC}/${NAME}_${INPROC}_event.root
  if [ ! -r $INFILE ]; then
    echo Input event file not found: $INFILE
    exit 4
  fi
fi 

runit
