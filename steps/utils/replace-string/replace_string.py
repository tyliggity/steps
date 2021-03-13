#!/usr/bin/env python3

import os
import json
import base64


def main():
    # required params
    original = os.environ["ORIGINAL_BASE64_STRING"]
    subStrSrc = os.environ["SUB_STR_SRC"]
    subStrDest = os.environ["SUB_STR_DEST"]

    decoded_content = base64.b64decode(original).decode("utf-8")
    result = base64.b64encode(
        decoded_content.replace(subStrSrc, subStrDest).encode("ascii")
    ).decode("ascii")

    output_object = {}
    output_object["output"] = result
    json_result = json.dumps(output_object)
    print("<-- END -->")
    print(json_result)


if __name__ == "__main__":
    main()
