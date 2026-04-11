import 'package:flutter/material.dart';
import 'package:retrovibed/uuidx.dart' as uuidx;

class Seed {
  static Seed global() => Seed(
    id: uuidx.min(),
    label: const Text("global"),
    description: const Text("public data"),
    tooltip: "global seed for public information, think creative commons",
    icon: Icons.public,
  );

  static Seed community(String id) => Seed(
    id: id,
    label: const Text("community"),
    description: const Text("Community shared content"),
    tooltip: "community seed, essentially public data",
    icon: Icons.group,
  );

  // personal isnt implemented yet.
  static Seed personal(String id) => Seed(
    id: id,
    label: const Text("personal"),
    description: const Text("Private to you"),
    tooltip:
        "your personal seed, used for information you want to keep private",
    icon: Icons.person,
  );

  static Seed unique(String id) => Seed(
    id: id,
    label: const Text("private"),
    description: const Text("Unique to this resource"),
    tooltip: "unique, not shared with anything else",
    icon: Icons.shield,
  );
  final String id;
  final Widget label;
  final Widget description;
  final String tooltip;
  final IconData icon;

  const Seed({
    required this.id,
    required this.label,
    required this.description,
    required this.icon,
    required this.tooltip,
  });
}

class Classifier {
  final String community;
  final String personal;

  const Classifier({required this.community, required this.personal});

  Seed classify(String seed) {
    if (seed == uuidx.min() || seed.isEmpty) {
      return Seed.global();
    }

    if (seed == community) {
      return Seed.community(seed);
    }

    if (seed == personal) {
      return Seed.personal(seed);
    }

    return Seed.unique(seed);
  }
}
