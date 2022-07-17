import os
import subprocess

tests_run = 0
tests_passed = 0

def run_test(file_path):
    global tests_run, tests_passed

    tests_run += 1

    expected_output = ""
    with open(file_path, "r", encoding="UTF-8") as file:
        expected_output = file.read()

    source_path = file_path.split(".")[0] + ".ava"
    result = subprocess.run(["go", "run", ".", "run", source_path], stdout=subprocess.PIPE)
    output = result.stdout.decode("utf-8")

    passed = expected_output == output

    if passed:
        tests_passed += 1
        print(f"{source_path}: PASS")
    else:
        print(f"{source_path}: FAIL")
        print(f"`{expected_output}` != `{output}`")

def main():
    for file in os.listdir("tests"):
        if file.endswith("txt"):
            run_test("tests/" + file)

    print()
    print(f"Total tests run: {tests_run}. ({tests_passed}/{tests_run})")

if __name__ == "__main__":
    main()