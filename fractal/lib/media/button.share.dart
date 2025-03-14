import 'package:uuid/uuid.dart' as uuid;
import 'package:flutter/material.dart';
import './media.pb.dart';

class ButtonShare extends StatelessWidget {
  final Media current;
  const ButtonShare({super.key, required this.current});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return IconButton(
      icon: Icon(Icons.share_outlined),
      color:
          uuid.UuidValue.fromString(current.torrentId) ==
                  uuid.Namespace.nil.uuidValue
              ? theme.disabledColor
              : const Color.fromARGB(255, 0, 255, 17),
      onPressed: () {
        print(
          "sharing management is not let implemented ${current.id} - ${current.torrentId}",
        );
      },
    );
  }
}
