import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'community.social.pb.dart';

class SocialEdit extends StatelessWidget {
  final CommunityPublisher compub;
  final void Function(CommunityPublisher) onChange;
  final bool autofocus;
  final bool readOnly;
  final List<Widget> actions;

  const SocialEdit({
    super.key,
    required this.compub,
    required this.onChange,
    this.autofocus = false,
    this.readOnly = false,
    this.actions = const [],
  });

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);

    return forms.Container(
      decoration: BoxDecoration(borderRadius: defaults.borderRadius),
      Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        spacing: defaults.spacing,
        children: [
          forms.Field(
            label: Text('Title Template'),
            trailing: actions,
            input: TextFormField(
              readOnly: readOnly,
              autofocus: !readOnly && autofocus,
              initialValue: compub.templateTitle,
              onChanged: (v) => onChange(compub..templateTitle = v.trim()),
              decoration: InputDecoration(
                hintText: "title template",
                border: OutlineInputBorder(),
              ),
            ),
          ),
          forms.Field(
            label: Text('Description Template'),
            input: TextFormField(
              initialValue: compub.templateDescription,
              onChanged: (v) => onChange(compub..templateDescription = v.trim()),
              decoration: InputDecoration(
                hintText: 'description template',
                border: OutlineInputBorder(),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
