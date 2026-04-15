import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;

class EditButton extends StatelessWidget {
  final VoidCallback onEdit;

  const EditButton({super.key, required this.onEdit});

  @override
  Widget build(BuildContext context) {
    return ds.LoadingIconButton(
      onPressed: () {
        onEdit();
        return Future.value();
      },
      icon: Icon(Icons.edit),
      tooltip: 'Edit Community',
      help: ds.Hint(const Text("edit community description and publish mode")),
    );
  }
}
