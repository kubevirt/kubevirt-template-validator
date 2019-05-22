#!/bin/bash
if [ -z "${V}" ]; then
	V=0
fi

RET=0
for testscript in $( ls ??-test-*.sh); do
	testname=$(basename -- "$testscript")
	testname="${testname%.*}"  # see http://www.gnu.org/software/bash/manual/html_node/Shell-Parameter-Expansion.html

	result="???"
	if [ "${V}" == "0" ]; then
		./$testscript &> /dev/null
	else
		printf "* TESTCASE [%-64s] START\n" $testscript
		./$testscript
	fi
	if [ "$?" == "0" ]; then
		result="OK"
	else
		if [ "$?" == "99" ] ; then
			result="SKIP"
		else
			result="FAILED"
			RET=1
		fi
	fi
	if [ "${V}" == "0" ]; then
		printf "* [%-64s] %s\n" $testscript $result
	else
		printf "  TESTCASE [%-64s] %s\n" $testscript $result
	fi
done
exit $RET
