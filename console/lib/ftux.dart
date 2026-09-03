import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'ftux/community.picker.dart';

export 'ftux/api.dart';

/// Shows the curated community-subscription picker once, the first time this
/// widget is mounted (i.e. right after login, when wrapped around the
/// authenticated app shell). Persists "seen" via [ds.Disclaimer]'s disk cache.
class AutoHelp extends StatefulWidget {
  final Widget child;
  const AutoHelp(this.child, {super.key});

  @override
  State<AutoHelp> createState() => _AutoHelpState();
}

class _AutoHelpState extends State<AutoHelp> {
  static const String _cacheid = 'ftux';

  @override
  Widget build(BuildContext context) {
    return ds.Disclaimer(
      widget.child,
      cacheid: _cacheid,
      overlay: ds.Masked(
        Center(
          child: SingleChildScrollView(
            child: CommunityPicker(
              onDone: () {
                ds.Disclaimer.acknowledge(_cacheid);
                setState(() {});
              },
            ),
          ),
        ),
      ),
    );
  }
}
