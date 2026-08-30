import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/design.kit/inputs.dart' as inputs;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/mimex.dart' as mimex;
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
            label: Text('URL'),
            input: TextFormField(
              readOnly: readOnly,
              autofocus: !readOnly && autofocus,
              initialValue: community.url,
              onChanged: (v) => onChange(community..url = v.trim()),
              decoration: InputDecoration(
                hintText: 'https://example.community.retrovibe.space',
                helper: Text(
                  community.url.isEmpty
                      ? 'leave blank to auto-generate a url'
                      : community.url,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                  style: TextStyle(color: Colors.grey, fontSize: 10),
                ),
                border: OutlineInputBorder(),
              ),
              validator: (value) {
                if (value == null || value.trim().isEmpty) {
                  return null;
                }
                final uri = Uri.tryParse(value.trim());
                if (uri == null || !uri.isAbsolute) {
                  return 'must be a valid url';
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
          forms.Field(
            label: Text('Mimetype'),
            input: inputs.Mimetype(
              value: community.mimetype,
              presets: [
                (label: Text('binary (default)') as Widget, value: mimex.binary),
                (label: Text('search plugin') as Widget, value: mimex.search),
                if (authn.developer(context).alpha) ...[
                  (label: Text('metadata archive') as Widget, value: mimex.mediaarchive),
                  (label: Text('neural') as Widget, value: mimex.neural),
                ],
              ],
              onChanged: (v) => onChange(community..mimetype = v),
            ),
            help: Text(
              'Mimetype is used to provide the default to content uploaded to this commmunity, generally speaking you should leave it defaulted',
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
