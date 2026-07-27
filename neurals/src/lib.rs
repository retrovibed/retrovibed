use std::ffi::{CStr, CString};
use std::os::raw::{c_char, c_int};
use std::sync::OnceLock;
use tract_onnx::prelude::*;

type Plan = SimplePlan<TypedFact, Box<dyn TypedOp>, Graph<TypedFact, Box<dyn TypedOp>>>;

struct Model {
    path: String,
    plan: Plan,
}

static INSTANCE: OnceLock<Result<Model, String>> = OnceLock::new();

fn run_predict(
    model_path: &str,
    input: &str,
    seq_len: usize,
    num_tokens: i64,
    pad: i64,
    bos: i64,
    eos: i64,
) -> Result<String, String> {
    let instance = INSTANCE.get_or_init(|| {
        let plan = tract_onnx::onnx()
            .model_for_path(model_path)
            .map_err(|e| format!("Failed to parse ONNX: {e}"))?
            .into_optimized()
            .map_err(|e| format!("Optimization failed: {e}"))?
            .into_runnable()
            .map_err(|e| format!("failed to build runnable plan: {e}"))?;
        Ok(Model {
            path: model_path.to_string(),
            plan,
        })
    });

    let model = instance.as_ref().map_err(|e| e.clone())?;

    let mut flat = vec![pad; seq_len];
    for (i, c) in input.chars().enumerate() {
        if i >= seq_len {
            break;
        }
        flat[i] = (c as i64) % num_tokens;
    }

    let src: Tensor = tract_ndarray::Array2::from_shape_vec((1, seq_len), flat)
        .map_err(|e| format!("{0} tensor error: {e}", model.path))?
        .into_dyn()
        .into();

    let mut outputs = model
        .plan
        .run(tvec!(src.into()))
        .map_err(|e| format!("{0} execution error: {e}", model.path))?;

    if outputs.is_empty() {
        return Err(format!("{0} model defines no outputs", model.path));
    }

    let logits = outputs
        .remove(0)
        .to_array_view::<f32>()
        .map_err(|e| format!("{0} output type mismatch: {e}", model.path))?
        .into_owned();

    let want_shape = [1, seq_len, num_tokens as usize];
    if logits.shape() != want_shape {
        return Err(format!(
            "{0} unexpected output shape: got {1:?}, want {want_shape:?}",
            model.path,
            logits.shape(),
        ));
    }

    let mut token_ids: Vec<i64> = Vec::new();
    for i in 0..seq_len {
        let tok = (0..num_tokens as usize)
            .max_by(|&a, &b| logits[[0, i, a]].partial_cmp(&logits[[0, i, b]]).unwrap())
            .unwrap_or(0) as i64;
        if tok == eos {
            break;
        }
        token_ids.push(tok);
    }

    Ok(token_ids
        .iter()
        .filter(|&&t| t != pad && t != bos && t != eos)
        .filter_map(|&t| char::from_u32(t as u32))
        .collect())
}

/// Predict output for input string using the ONNX model at model_path.
/// Writes a null-terminated result into output (up to output_len bytes).
/// Returns 0 on success, 1 on error.
#[no_mangle]
pub extern "C" fn predict(
    model_path: *const c_char,
    input: *const c_char,
    seq_len: usize,
    num_tokens: i64,
    pad: i64,
    bos: i64,
    eos: i64,
    output: *mut c_char,
    output_len: usize,
) -> c_int {
    let model_path = unsafe { CStr::from_ptr(model_path) }.to_string_lossy();
    let input = unsafe { CStr::from_ptr(input) }.to_string_lossy();

    match run_predict(&model_path, &input, seq_len, num_tokens, pad, bos, eos) {
        Ok(result) => {
            let s = CString::new(result).unwrap_or_default();
            let bytes = s.as_bytes_with_nul();
            let copy_len = bytes.len().min(output_len);
            unsafe {
                std::ptr::copy_nonoverlapping(bytes.as_ptr() as *const c_char, output, copy_len)
            };
            0
        }
        // The C ABI only carries a 0/1 status back to Go, so the failure
        // reason (which run_predict branch, which model, etc.) would
        // otherwise be silently discarded. Logging it here to stderr is
        // the only way to see why a specific predict call failed.
        Err(e) => {
            eprintln!("predict error: {e}");
            1
        }
    }
}
