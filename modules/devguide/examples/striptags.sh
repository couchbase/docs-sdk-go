#!/bin/sh
here="`pwd`"
cd "$(dirname "$0")"
if [ -z "$gitmod" ] ; then
	gitmod=0
fi
target="$(cd ../../../../devguide-examples ; pwd)" # only 3 ../ if using submodule
# exclude striptags.sh and non-git files like nodejs/node_modules
find * |\
egrep -v '^striptags.sh|^nodejs/node_modules|^nodejs/package-lock.json' |\
 while read f ; do
  if [ -d $f ] ; then
    if [ ! -d $target/$f ] ; then
        echo Creating directory $target/$f
	mkdir $target/$f || exit
	if [ $gitmod -eq 1 ] ; then 
		echo git add $target/$f
		( cd $target ; git add $f )
	fi
    fi
  else
	egrep -v '#tag|#end' $f > /tmp/stripped.$$
        test  $? -gt 1 && exit
	diff -q /tmp/stripped.$$ $target/$f 
        if [ $? -ne 0 ] ; then 
		mv /tmp/stripped.$$ $target/$f || exit
		if [ $gitmod -eq 1 ] ; then 
			echo git add $target/$f
			( cd $target ; git add $f )
		fi
	fi
  fi
done
rm -f /tmp/stripped.$$
cd "$here"
