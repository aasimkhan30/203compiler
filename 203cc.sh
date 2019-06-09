#!/bin/bash
# ------------------------------------------------------------------
# [Author] Title
#          Description
# ------------------------------------------------------------------

VERSION=0.1.0
SUBJECT=some-unique-id
USAGE="Usage: 203cc.sh -l <go|python> filename.c"

# --- Options processing -------------------------------------------
if [ $# == 0 ] ; then
    echo $USAGE
    exit 1;
fi

while getopts ":l:vh" optname
  do
    case "$optname" in
      "l")
      echo "Language selected for compilation : $OPTARG"
      case "$OPTARG" in
       "go")
        go run ./go_lang/gocc.go $3
        ;;
        "python")
        python ./python/pycc.py $3
        ;;
        *)
        echo "Check your language parameter"
      esac
        ;;
      "?")
        echo "Unknown option $OPTARG"
        exit 0;
        ;;
      ":")
        echo "No argument value for option $OPTARG"
        exit 0;
        ;;
      *)
        echo "Unknown error while processing options"
        exit 0;
        ;;
    esac
  done

shift $(($OPTIND - 1))
