import argparse
import json

if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Create Foundational Keys List from Internal Record."
    )
    parser.add_argument(
        "-json-file",
        default="keys.json",
        dest="json_file",
        help="tab seperate ecdsa and bls keys",
        type=str,
    )
    parser.add_argument(
        "-out",
        default="accounts-generated-go.txt",
        dest="out_file",
        help="file compatible with foundational go",
        type=str,
    )
    parser.add_argument(
        "-index",
        default=0,
        dest="index",
        help="index of where you want to start from",
        type=int,
    )
    args = parser.parse_args()
    g = open(args.json_file, "r")
    f = open(args.out_file, "w")
    d = json.load(g)
    index = args.index
    for item in d:
        string = (
            "{Index: "
            + '"'
            + str(index)
            + '"'
            + ","
            + " "
            + "ShardID: "
            + str(item['shard-id'])
            + ","
            + " "
            + "Address: "
            + '"'
            + item['hex-address']
            + '"'
            + ","
            + " "
            + "BLSPublicKey: "
            + '"'
            + item['public-key']
            + '"'
            + "}"
            + ","
        )
        f.write(string + "\n")
        index = index + 1
    g.close()
