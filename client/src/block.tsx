import React, { useState } from "react";
import { Theme, createStyles, makeStyles } from "@material-ui/core/styles";
import Paper from "@material-ui/core/Paper";
import { Typography, Box } from "@material-ui/core";

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
      width: "100%",
      display: "flex",
      flexWrap: "wrap",
      "& > *": {
        margin: theme.spacing(1),
        width: theme.spacing(16),
        height: theme.spacing(16),
      },
    },
    paper: {
      width: "40%",
      margin: "10px",
      height: "100%",
      border: '1px solid purple'
    },
    textFormat: {
      fontSize: "14px",
    },
  })
);

export const BlockComponent = (props) => {
  const classes = useStyles();
  const blockHighlight = props.block.Data === "GENESIS BLOCK" 
    ? {border: '2px solid green'} 
    : {border: '1px solid black'} 

  return (
    <Box border={blockHighlight} className={classes.root}>
      <Paper variant="outlined" className={classes.paper} style={blockHighlight} elevation={4}>
        <Box margin="5px">
          <Typography
            className={classes.textFormat}
            variant="h4"
          >{`ID: ${props.block.ID}`}</Typography>
        </Box>
        <Box margin="5px">
          <Typography
            className={classes.textFormat}
            variant="h6"
          >{`Index: ${props.block.Index}`}</Typography>
        </Box>
        <Box margin="5px">
          <Typography
            className={classes.textFormat}
            variant="h4"
          >{`Hash: ${props.block.Hash}`}</Typography>
        </Box>
        <Box margin="5px">
          <Typography
            className={classes.textFormat}
            variant="h4"
          >{`PreviousHash: ${props.block.PreviousHash}`}</Typography>
        </Box>
        <Box margin="5px">
          <Typography
            className={classes.textFormat}
            variant="h4"
          >{`Timestamp: ${props.block.Timestamp}`}</Typography>
        </Box>
        <Box margin="5px">
          <Typography
            className={classes.textFormat}
            variant="h4"
          >{`Data: ${props.block.Data}`}</Typography>
        </Box>
      </Paper>
    </Box>
  );
};
