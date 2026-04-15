import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/community/api.dart';
import 'package:retrovibed/authn.dart' as authn;

class DeleteButton extends StatelessWidget {
  final Community community;
  final void Function(Community)? onDeleted;

  const DeleteButton({super.key, required this.community, this.onDeleted});

  Future<void> _confirm(BuildContext context) {
    return showDialog<bool>(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: Text('Delete Community'),
          content: Text(
            'Are you sure you want to delete "${community.domain}"? This action cannot be undone.',
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(false),
              child: Text('Cancel'),
            ),
            TextButton(
              onPressed: () => Navigator.of(context).pop(true),
              style: TextButton.styleFrom(foregroundColor: Colors.red),
              child: Text('Delete'),
            ),
          ],
        );
      },
    ).then((confirmed) {
      if (confirmed != true) return Future.value();
      final auth = [authn.DeeppoolAuthzCache.bearer(context)];
      return httpx.withRetry(() => API.delete(community.id, options: auth)).then((v) => onDeleted?.call(v.community));
    });
  }

  @override
  Widget build(BuildContext context) {
    return ds.LoadingIconButton(
      onPressed: () => _confirm(context),
      icon: Icon(Icons.delete, color: Colors.red),
      tooltip: 'Delete Community',
      help: ds.Hint(const Text("permanently remove this community")),
    );
  }
}
