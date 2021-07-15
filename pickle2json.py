'''
Convert a pkl file into json file
'''
import sys
import os
import pickle
import json
import gzip

class JsonEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, set):
            return { k: {} for k in obj }
        return json.JSONEncoder.default(self, obj)

def convert_dict_to_json(file_path, output_path):
    with gzip.open(file_path, 'rb') as fpkl, open(output_path, 'w') as fjson:
        weights, classes = pickle.load(fpkl)
        json.dump({'Weights': weights, 'Classes': classes }, fjson, cls=JsonEncoder, ensure_ascii=True, sort_keys=True)


def main():
    if sys.argv[1] and os.path.isfile(sys.argv[1]):
        file_path = sys.argv[1]
        output_path = './' + os.path.basename(file_path)
        if sys.argv[2]:
            output_path = sys.argv[2]
        print("Processing %s to %s ..." % (file_path, output_path))
        convert_dict_to_json(file_path, output_path)
    else:
        print("Usage: %s abs_file_path" % (__file__))


if __name__ == '__main__':
    main()
