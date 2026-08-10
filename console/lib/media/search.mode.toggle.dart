import 'package:flutter/material.dart';

// Reused by both the library and discovery search trays to flip between
// searching the local library and discovering content over the network.
PopupMenuItem<String> SearchModeToggle({required bool discovering, required VoidCallback onToggle}) {
  return PopupMenuItem<String>(
    onTap: onToggle,
    child: ListTile(
      leading: Icon(discovering ? Icons.check : Icons.travel_explore),
      title: const Text("Search"),
    ),
  );
}
