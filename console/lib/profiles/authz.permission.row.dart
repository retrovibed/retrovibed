import 'package:flutter/material.dart';
import 'package:retrovibed/design.kit/forms.dart' as forms;

class AuthzPermissionRow extends StatelessWidget {
  final String label;
  final String description;
  final bool value;
  final void Function(bool)? onChanged;

  const AuthzPermissionRow(
    this.label, {
    super.key,
    required this.description,
    this.value = false,
    this.onChanged,
  });

  @override
  Widget build(BuildContext context) {
    return ConstrainedBox(
      constraints: const BoxConstraints(maxWidth: 288.0),
      child: forms.Checkbox(
        Text(label),
        description: Text(
          description,
          maxLines: 2,
          overflow: TextOverflow.ellipsis,
        ),
        value: value,
        onChanged: onChanged != null ? (v) => onChanged!(v ?? false) : null,
      ),
    );
  }
}
