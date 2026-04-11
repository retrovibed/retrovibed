import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/community/api.dart';

class ShareButton extends StatelessWidget {
  final Community community;

  const ShareButton({super.key, required this.community});

  @override
  Widget build(BuildContext context) {
    return ds.LoadingIconButton(
      onPressed: () => Clipboard.setData(ClipboardData(text: community.url)),
      icon: Icon(Icons.share),
      tooltip: 'Share',
      help: ds.Hint(
        label: const Text("Share"),
        description: const Text("copies the community link to your clipboard"),
      ),
    );
  }
}
