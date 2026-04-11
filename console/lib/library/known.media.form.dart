import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/media/media.known.pb.dart';

class KnownMediaForm extends StatefulWidget {
	final void Function(KnownCreateRequest) onChange;
	final KnownCreateRequest? initial;

	const KnownMediaForm({
		super.key,
		required this.onChange,
		this.initial,
	});

	@override
	State<KnownMediaForm> createState() => _KnownMediaFormState();
}

class _KnownMediaFormState extends State<KnownMediaForm> {
	late KnownCreateRequest _current;

	@override
	void initState() {
		super.initState();
		_current = widget.initial ?? KnownCreateRequest(
			known: Known(
				released: DateTime.now().toIso8601String().split('T').first,
				adult: false,
			),
		);
	}

	void _update(void Function(Known) fn) {
		fn(_current.known);
		widget.onChange(_current);
		setState(() {});
	}

	@override
	Widget build(BuildContext context) {
		final defaults = ds.Defaults.of(context);
		final known = _current.known;

		return forms.Container(
			decoration: BoxDecoration(
				borderRadius: defaults.borderRadius,
			),
			Column(
				mainAxisSize: MainAxisSize.min,
				spacing: defaults.spacing,
				children: [
					forms.Field(
						label: Text('Title'),
						input: TextFormField(
							initialValue: known.description,
							decoration: InputDecoration(
								hintText: 'Content title',
								border: OutlineInputBorder(),
							),
							onChanged: (v) => _update((k) => k.description = v),
						),
					),
					forms.Field(
						label: Text('Summary'),
						input: TextFormField(
							initialValue: known.summary,
							decoration: InputDecoration(
								hintText: 'Content summary',
								border: OutlineInputBorder(),
							),
							maxLines: 3,
							onChanged: (v) => _update((k) => k.summary = v),
						),
					),
					forms.Field(
						label: Text('Image URL (optional)'),
						input: TextFormField(
							initialValue: known.image,
							decoration: InputDecoration(
								hintText: 'https://example.com/image.jpg',
								border: OutlineInputBorder(),
							),
							onChanged: (v) => _update((k) => k.image = v),
						),
					),
					forms.Field(
						label: Text('Release Date'),
						input: TextFormField(
							initialValue: known.released,
							decoration: InputDecoration(
								hintText: 'YYYY-MM-DD',
								border: OutlineInputBorder(),
							),
							onChanged: (v) => _update((k) => k.released = v),
						),
					),
					forms.Checkbox(
						Text('Adult content'),
						value: known.adult,
						onChanged: (v) => _update((k) => k.adult = v ?? false),
						description: Text('Mark if this content is for adults only'),
					),
				],
			),
		);
	}
}
