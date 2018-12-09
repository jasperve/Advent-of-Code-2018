import networkx as nx

G = nx.DiGraph()
with open('input.txt') as f:
    for line in f:
        instruction = line.strip().split()
        G.add_edge(instruction[1], instruction[7])

tasks = {}
seconds = 0
available = True

while len(G) or len(tasks):

    selectedKeys = list()

    for k, v in tasks.items():
        if v == seconds:
            selectedKeys.append(k)

    for k in selectedKeys:
        if k in tasks:
            print("finished with node: ", node)
            del tasks[k]
            G.remove_node(k)

    nodes = [v for v, d in G.in_degree() if d == 0]

    for node in nodes:
        if node not in tasks and len(tasks) < 5:
            print("adding: ", node)
            tasks[node] = ord(node) + seconds - 4

    if not len(tasks):
        break
    else:
        seconds = min(tasks.values())

print(seconds)