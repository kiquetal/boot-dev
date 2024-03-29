from main import Graph

run_cases = [
    (
        (
            ["Gondor", "Rivendell"],
            ["Gondor", "Lothlorien"],
            ["Rivendell", "Lothlorien"],
            ["Lothlorien", "Minas Tirith"],
            ["Rivendell", "Minas Tirith"],
            ["Rivendell", "Bree"],
            ["Bree", "Minas Tirith"],
        ),
        {
            "Gondor": {"Rivendell", "Lothlorien"},
            "Rivendell": {"Gondor", "Bree", "Minas Tirith", "Lothlorien"},
            "Bree": {"Rivendell", "Minas Tirith"},
            "Minas Tirith": {"Rivendell", "Bree", "Lothlorien"},
            "Lothlorien": {"Gondor", "Rivendell", "Minas Tirith"},
        },
    ),
    (
        (
            ["Rivendell", "Bree"],
            ["Bree", "Minas Tirith"],
            ["Minas Tirith", "Lothlorien"],
        ),
        {
            "Rivendell": {"Bree"},
            "Bree": {"Rivendell", "Minas Tirith"},
            "Minas Tirith": {"Bree", "Lothlorien"},
            "Lothlorien": {"Minas Tirith"},
        },
    ),
]
submit_cases = run_cases + [
    ((), {}),
    ((["Gondor", "Rivendell"],), {"Gondor": {"Rivendell"}, "Rivendell": {"Gondor"}}),
    ((["Rivendell", "Gondor"],), {"Gondor": {"Rivendell"}, "Rivendell": {"Gondor"}}),
    (
        (
            ["Bree", "Minas Tirith"],
            ["Minas Tirith", "Lothlorien"],
            ["Bree", "Lothlorien"],
            ["Rivendell", "Bree"],
        ),
        {
            "Rivendell": {"Bree"},
            "Bree": {"Rivendell", "Minas Tirith", "Lothlorien"},
            "Minas Tirith": {"Bree", "Lothlorien"},
            "Lothlorien": {"Bree", "Minas Tirith"},
        },
    ),
]


def test(edges, expected):
    try:
        print("---------------------------------")
        print("Creating adjacency list from edges:")
        print(f"{edges}\n")
        actual = Graph()
        for edge in edges:
            actual.add_edge(*edge)
        print("Expecting graph:")
        for k, v in expected.items():
            print(f"{k}: {v}")
        print("\nActual graph:")
        for k, v in actual.graph.items():
            print(f"{k}: {v}")
        if actual.graph == expected:
            print("\nPass")
            return True
        print("\nFail")
        return False
    except Exception as e:
        print("\nFail")
        print(e)
        return False


def main():
    passed = 0
    failed = 0
    for test_case in test_cases:
        correct = test(*test_case)
        if correct:
            passed += 1
        else:
            failed += 1
    if failed == 0:
        print("============= PASS ==============")
    else:
        print("============= FAIL ==============")
    print(f"{passed} passed, {failed} failed")


test_cases = submit_cases
if "__RUN__" in globals():
    test_cases = run_cases

main()
