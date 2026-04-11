import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/community/community.pb.dart';
import 'package:retrovibed/community/publish.mode.edit.dart';
import 'package:retrovibed/community/visibility.selector.dart';

class CommunityEdit extends StatelessWidget {
  final Community community;
  final void Function(Community) onChange;
  final bool autofocus;
  final bool readOnly;

  const CommunityEdit({
    super.key,
    required this.community,
    required this.onChange,
    this.autofocus = false,
    this.readOnly = false,
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
            label: Text('Domain'),
            input: TextFormField(
              readOnly: readOnly,
              autofocus: !readOnly && autofocus,
              initialValue: community.domain,
              onChanged: (v) => onChange(community..domain = v.trim()),
              decoration: InputDecoration(
                hintText: 'example',
                helper: Text(
                  'https://${community.domain.isEmpty ? 'example' : community.domain}.community.retrovibe.space',
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                  style: TextStyle(color: Colors.grey, fontSize: 10),
                ),
                border: OutlineInputBorder(),
              ),
              validator: (value) {
                if (value == null || value.trim().isEmpty) {
                  return 'Domain is required';
                }
                return null;
              },
            ),
          ),
          forms.Field(
            label: Text('Description'),
            input: TextFormField(
              autofocus: readOnly && autofocus,
              initialValue: community.description,
              onChanged: (v) => onChange(community..description = v.trim()),
              decoration: InputDecoration(
                hintText: 'description',
                border: OutlineInputBorder(),
              ),
            ),
          ),
          VisibilitySelector(
            hidden: community.hidden,
            onChanged: (v) => onChange(community..hidden = v),
          ),
          PublishModeEdit(
            publishMode: community.defaultPublishMode,
            onChanged: (v) => onChange(community..defaultPublishMode = v),
          ),
        ],
      ),
    );
  }
}
