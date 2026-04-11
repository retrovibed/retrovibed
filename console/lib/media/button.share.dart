import 'package:uuid/uuid.dart' as uuid;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import './media.pb.dart';

class ButtonShare extends StatelessWidget {
  final Media current;
  const ButtonShare({super.key, required this.current});

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final theme = Theme.of(context);

    return IconButton(
      icon: Icon(Icons.share_outlined),
      color:
          uuidx.fromString(current.torrentId) == uuid.Namespace.nil.uuidValue
              ? theme.disabledColor
              : defaults.success,
      onPressed: () {
        print(
          "sharing management is not let implemented ${current.id} - ${current.torrentId}",
        );
      },
    );
  }
}
