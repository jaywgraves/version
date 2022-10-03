FILE=$1
FNAME="${FILE%.*}"
FEXT="${FILE#*.}"
FDATE=$(TZ=US/Central date -r $FILE "+%Y%m%d.%H%M%S")
NEWFILE=$FNAME.$FDATE.$FEXT
echo "copying " $FILE "to" $NEWFILE
cp $FILE $NEWFILE