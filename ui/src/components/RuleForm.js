import React, { Component } from 'react';

import { withStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import RadioGroup from '@material-ui/core/RadioGroup';
import Radio from '@material-ui/core/Radio';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import axios from "axios";

const styles = theme => ({
  table: {
    minWidth: 320
  }
});

class RuleForm extends Component {
  state = {
    open: false,
    rule_use: "true",
    rule_alarm_min: "",
    rule_alarm_max: "",
  };

  setOpen = (val) => {
    this.setState({
      open: val
    })
  }

  render() {
    const { open, rule_use, rule_alarm_min, rule_alarm_max } = this.state;

    const openWindow = (val) => {
      this.setOpen(true);
    }

    const closeWindow = () => {
      this.setOpen(false);
    }

    const handleChange = (e) => {
      this.setState({
        rule_use: e.target.value,
      });
    }

    const setMin = (e) => {
      this.setState({
        rule_alarm_min: e.target.value,
      });
    }

    const setMax = (e) => {
      this.setState({
        rule_alarm_max: e.target.value,
      });
    }

    const runPost = async () => {
      const baseurl = (process.env.NODE_ENV === "development") ? "http://localhost:8080":"";

      const reqJson = {
        "Use": JSON.parse(rule_use),
        "AlarmMin": Number(rule_alarm_min),
        "AlarmMax": Number(rule_alarm_max),
      }
      try {
        const res = await axios.post(`${baseurl}/api/rule`, reqJson);
        console.log(res);
        closeWindow();
      } catch (err) {
          alert(`Error: ${err}`);
      }
    }

    return (
      <div>
        <button onClick={openWindow}>SetRule</button>
        <Dialog open={open} onClose={closeWindow} aria-labelledby="form-dialog-title">
          <DialogTitle id="form-dialog-title">SetRule</DialogTitle>
          <DialogContent>
            <RadioGroup value={rule_use} onChange={handleChange}>
              <FormControlLabel value="true" control={<Radio />} label="True" />
              <FormControlLabel value="false" control={<Radio />} label="False" />
            </RadioGroup>
            <TextField
              autoFocus
              margin="dense"
              id="alarmMin"
              label="Alarm Min"
              type="text"
              value={rule_alarm_min}
              onChange={setMin}
              fullWidth
            />
            <TextField
              margin="dense"
              id="alarmMin"
              label="Alarm Max"
              type="text"
              value={rule_alarm_max}
              onChange={setMax}
              fullWidth
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={runPost} color="primary">
              Submit
            </Button>
            <Button onClick={closeWindow} color="primary">
              Cancel
            </Button>
          </DialogActions>
        </Dialog>
      </div>
    );
  }
}

export default withStyles(styles)(RuleForm);
