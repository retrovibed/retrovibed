import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/community/api.dart';
import 'package:retrovibed/community/publish.container.dart';

class PublishButton extends StatelessWidget {
  final Community community;

  const PublishButton({super.key, required this.community});

  @override
  Widget build(BuildContext context) {
    return ds.LoadingIconButton(
      onPressed: () => ds.modals.asyncfn<void>(
        context,
        (completion) => PublishContainer(
          onPublished: completion.complete,
          onCancel: completion.complete,
          community: community,
        ),
      ),
      icon: Icon(Icons.publish),
      tooltip: 'Publish Content',
      help: ds.Hint(const Text("share library content to this community")),
    );
  }
}
