#!/bin/bash
# Copyright Contributors to the Open Cluster Management project

# TESTED ON MAC!

# NOTE: When running against a node repo, delete the node_modules directories first!  Then npm ci once all the
#       copyright changes are incorporated.

# set -x
TMP_FILE="tmp_file"

IGNORE_COPYRIGHT_FILES="
vendor
client
_generated.deepcopy.go
crd.yaml
.github
.gitignore
CHANGELOG
"

verify="${VERIFY:-}"

ALL_FILES=$(git ls-files | 
 grep -v -f <(echo "$IGNORE_COPYRIGHT_FILES" | grep . | sed 's/\([.|]\)/\1/g; s/\?/./g ; s/\*/.*/g'))

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..

COMMUNITY_COPY_HEADER_FILE="${SCRIPT_ROOT}/hack/boilerplate.txt"

if [ ! -f $COMMUNITY_COPY_HEADER_FILE ]; then
  echo "File $COMMUNITY_COPY_HEADER_FILE not found!"
  exit 1
fi

COMMUNITY_COPY_HEADER_STRING=$(cat $COMMUNITY_COPY_HEADER_FILE | sed 's#^// ##')

# NOTE: Only use one newline or javascript and typescript linter/prettier will complain about the extra blank lines
NEWLINE="\n"

ERROR=false

for FILE in $ALL_FILES
do
    if [[ -d $FILE ]] ; then
        continue
    fi

    COMMENT_START="# "
    COMMENT_END=""

    if [[ $FILE  == *".go" ]]; then
        COMMENT_START="// "
    fi

    if [[ $FILE  == *".ts" || $FILE  == *".tsx" || $FILE  == *".js" ]]; then
        COMMENT_START="/* "
        COMMENT_END=" */"
    fi

    if [[ $FILE  == *".md" ]]; then
        COMMENT_START="\[comment\]: # ( "
        COMMENT_END=" )"
    fi

    if [[ $FILE  == *".html" ]]; then
        COMMENT_START="<!-- "
        COMMENT_END=" -->"
    fi

    if [[ $FILE  == *".go"       \
            || $FILE == *".yaml" \
            || $FILE == *".yml"  \
            || $FILE == *".sh"   \
            || $FILE == *".js"   \
            || $FILE == *".ts"   \
            || $FILE == *".tsx"   \
            || $FILE == *"Dockerfile" \
            || $FILE == *"Makefile"  \
            || $FILE == *".mk" \
            || $FILE == *"Dockerfile.prow" \
            || $FILE == *"Makefile.prow"  \
            || $FILE == *".gitignore"  \
            || $FILE == *".md"  ]]; then

        COMMUNITY_HEADER_AS_COMMENT="$COMMENT_START$COMMUNITY_COPY_HEADER_STRING$COMMENT_END"

        if [ -f "$FILE" ] && ! grep -q "$COMMUNITY_COPY_HEADER_STRING" "$FILE"; then
          if [ "$verify" = true ]; then
            echo "$FILE"
            ERROR=true
          else
            echo "Adding copyright to $FILE"
                        (cat <<EOM
$COMMUNITY_HEADER_AS_COMMENT
$(cat $FILE)
EOM
) > "$TMP_FILE" && mv "$TMP_FILE" "$FILE"
          fi
        fi
    fi
done

if [ $ERROR = true ]
then
  exit 1
fi
rm -f $TMP_FILE
